require 'fog'

module Deployer
  unless ENV['AWS_ACCESS_KEY_ID'] && ENV['AWS_SECRET_ACCESS_KEY']
    raise "Need credentials at $AWS_ACCESS_KEY_ID and $AWS_SECRET_ACCESS_KEY"
  end

  unless ENV['PDFCOMBINER_PEM']
    raise "Please set $PDFCOMBINER_PEM to the path to the instance private ssh key path"
  end

  raise "No file at $PDFCOMBINER_PEM" unless File.exists?(ENV['PDFCOMBINER_PEM'])

  STACK_NAME = 'pdfcombiner'
  CREDS = { :aws_access_key_id => ENV['AWS_ACCESS_KEY_ID'],
            :aws_secret_access_key => ENV['AWS_SECRET_ACCESS_KEY'] }

  # Update the Cloudformation stack with the given parameters.
  # By default, only update the DeployedAt timstamp, which will not actually
  # do anything until new instances are spun up or cfn-init is run on
  # existing instances.  If you pass in InstanceType (or any other attribute
  # change that requires a restart), it will trigger a rolling rebuild of all
  # instances with the new size.
  def self.update_stack!(params_to_update={})
    raise 'aborting deploy' unless safe_to_update?
    new_params = deploy_timestamp.merge(params_to_update)
    CloudFormation.update_stack(new_params)
    logger.info("Updated stack with params: #{new_params}")
  end

  # Runs the cfn-init command on any existing instances via ssh, which will
  # copy the current binary and restart the service.
  def self.copy_binary_and_restart!
    asg_id = CloudFormation.autoscaling_group_id
    asg = AutoScalingGroup.new(asg_id)
    deploy_instances = asg.running_instance_ids
    SshTool.new(deploy_instances).redeploy
    logger.info("Redeployed and restarted service on #{deploy_instances}")
  end

  # Resize all servers in the scaling group to the given instance size.
  # TODO provide estimate based on server count and current delay param
  def self.resize!(new_size)
    if self.update_stack!('InstanceType' => new_size)
      logger.info("Finished setting new size to '#{new_size}'.  Rolling restart in progress")
      logger.info(progress_message)
    end
  end

  private

  def self.safe_to_update?
    case status = CloudFormation.stack['StackStatus']
    when /COMPLETE$/
      true
    else
      logger.error("Can't update stack while in state '#{status}'.")
      logger.error(progress_message)
      false
    end
  end

  def self.logger
    @logger ||= Logger.new($stdout)
  end

  def self.progress_message
    url="https://console.aws.amazon.com/cloudformation/home?region=us-east-1"+
    "#ConsoleState:viewing=ACTIVE&stack=#{STACK_NAME}&tab=Events"
    "Check #{url} to monitor progress"
  end

  def self.deploy_timestamp
    {"DeployString" => "Deployed by #{ENV['USER']} at #{Time.now}"}
  end

  class CloudFormation
    class << self

      # Update the stack with the current template file and any new parameters
      def update_stack(updated_params)
        new_params = current_params.merge(updated_params)
        fog.update_stack STACK_NAME,
          "TemplateBody" => template_contents,
          "Parameters" => new_params
      end

      # Examine the stack to determine the associated autoscaling group id.
      def autoscaling_group_id
        asg = stack_resources.detect{|x| x['ResourceType'] == "AWS::AutoScaling::AutoScalingGroup"}
        asg && asg['PhysicalResourceId']
      end

      # The currently deployed parameters of the stack.
      def current_params
        stack["Parameters"].inject({}, &extract_params)
      end

      def stack
        @stack ||= stack_with_name(STACK_NAME)
      end

      private

      def extract_params
        ->(memo, curr) {
          memo.merge( curr['ParameterKey'] => curr['ParameterValue'] )
        }
      end

      def stack_resources
        response = fog.describe_stack_resources('StackName' => STACK_NAME).body
        response['StackResources']
      end

      def template_contents
        pwd = File.expand_path File.dirname(__FILE__)
        File.read( File.join(pwd, '../cloudformation.json'))
      end

      def stack_with_name(name)
        stacks.detect{|x| x['StackName'] == name} or raise "stack #{name} not found"
      end

      def stacks
        @stacks ||= fog.describe_stacks.body['Stacks']
      end

      def fog
        @fog ||= Fog::AWS::CloudFormation.new(CREDS)
      end

    end
  end

  class AutoScalingGroup
    attr_accessor :group

    def initialize(group_id)
      @group = find(group_id) or raise "No such group id: '#{group_id}'"
    end

    def running_instance_ids
      group.instances.map(&:id)
    end

    private

    def find(group_id)
      fog.groups.detect{|x| x.id == group_id }
    end

    def fog
      @fog ||= Fog::AWS::AutoScaling.new(CREDS)
    end

  end

  class SshTool
    REDEPLOY_COMMAND = "sudo cfn-init -v -s #{STACK_NAME} -r LaunchConfig -c ALL"
    RESTART_COMMAND = 'sudo service pdfcombiner restart'

    def initialize(instance_ids)
      @ec2 = Fog::Compute.new(CREDS.merge(provider: 'AWS'))
      @instance_ids = Array(instance_ids)
      @running_instances = find_instances_by_ids(instance_ids)
    end

    # TODO check for service status after/during
    # handle failed ssh
    def redeploy
      Deployer.logger.info("About to deploy to #{friendly_server_names}")
      @running_instances.each do |instance|
        instance.private_key_path = ENV['PDFCOMBINER_PEM']
        responses = instance.ssh([RESTART_COMMAND, REDEPLOY_COMMAND])
        if responses.any?{ |response| response.status != 0 }
          raise "redeploy failed on #{instance}: #{responses}"
        end
      end
    end

    def find_instances_by_ids(ids)
      @ec2.servers.select(&in_desired_set?).tap(&check_empty)
    end

    def in_desired_set?
      ->(server) {
        @instance_ids.include?(server.id) && server.state == 'running'
      }
    end

    def check_empty
      ->(servers){ raise "no running servers to deploy to!  Tried #{@instance_ids}" if servers.empty? }
    end

    def friendly_server_names
      @running_instances.map{|s| "#{s.id}(#{s.dns_name})"}.join(', ')
    end

  end
end

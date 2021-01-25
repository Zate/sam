import logging
import sys
import os
from typing import Any
try:
    sys.path.insert(0, os.path.join(os.path.dirname(__file__), "lib"))
except Exception as e:
    sys.exit("Cannot load lib folder " + str(e))
    pass

from splunklib.modularinput import Argument, EventWriter, InputDefinition, Scheme, Script, ValidationDefinition  # type: ignore

from validation import Validation       # this is our argument validation
# # this loads the events from the API
from event_loader import EventLoader


logging.root.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(levelname)s %(message)s')
handler = logging.StreamHandler(stream=sys.stderr)
handler.setFormatter(formatter)
logging.root.addHandler(handler)


class TerraformCloudModInput(Script):  # type: ignore
    def get_scheme(self) -> Scheme:
        """Loads up a scheme per the Splunk docs
        https://docs.splunk.com/Documentation/Splunk/8.0.3/AdvancedDev/ModInputsScripts#Define_a_scheme_for_introspection
        https://dev.splunk.com/enterprise/docs/python/sdk-python/howtousesplunkpython/howtocreatemodpy/#The-get_scheme-method
        """
        # logging.debug("Processing scheme request")
        scheme = Scheme("Terraform Cloud for Splunk")
        scheme.description = "Gets data from Terraform"
        scheme.streaming_mode = "xml"
        scheme.use_external_validation = True

        organization_token = Argument(
            name="organization_token",
            description="Your Terraform Organization Token",
            data_type=Argument.data_type_string,
            required_on_edit=True,
            required_on_create=True,
            title="Terraform Organization Token"
        )
        scheme.add_argument(organization_token)

        api_url = Argument(
            name="api_url",
            description="Your Terraform URL",
            data_type=Argument.data_type_string,
            required_on_edit=True,
            required_on_create=True,
            title="Terraform URL"
        )
        scheme.add_argument(api_url)
        return scheme

    def validate_input(self, validation_definition: ValidationDefinition) -> Any:
        """ Validates the input by calling home
        https://docs.splunk.com/Documentation/Splunk/8.0.3/AdvancedDev/ModInputsValidate
        https://dev.splunk.com/enterprise/docs/python/sdk-python/howtousesplunkpython/howtocreatemodpy/#The-validate_input-method
        """

        # logging.debug("Processing validate argument request")
        Validation.validate_arguments(validation_definition)

    def stream_events(self,	input_definition: InputDefinition, ew: EventWriter) -> None:
        """Streams events from the API into XML format
        https://dev.splunk.com/enterprise/docs/python/sdk-python/howtousesplunkpython/howtocreatemodpy/#The-stream_events-method
        """
        # logging.debug("Processing stream events request")
        stanza = list(input_definition.inputs.keys())[0]
        host = input_definition.inputs[stanza]["api_url"]
        trails = EventLoader.load_events(input_definition, ew.log, stanza)
        for trail in trails:
            ew.write_event(trail.to_event(stanza, host))


if __name__ == '__main__':
    sys.exit(TerraformCloudModInput().run(sys.argv))

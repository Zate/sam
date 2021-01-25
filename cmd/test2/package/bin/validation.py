import sys
import logging
from typing import Any, Dict
from splunklib.modularinput import ValidationDefinition  # type: ignore
from error_formatter import ErrorFormatter
from tf_client import TFClient


class ValidationException(Exception):
    pass


class Validation:
    """Handles calling TF to make sure this organization token matches audit logging entitlement"""
    @classmethod
    def validate_arguments(self, validation_definition: ValidationDefinition) -> bool:
        """
        This takes in the passed in validation definition and verifies it.
        Then it calls the organization entitlements API endpoint
        If they have access to audit_logging, we continue, otherwise raise an exception

        Args:
            validation_definition: config from splunk

        Returns:
            True if they have entitlement
        """
        try:
            # loads and validates the configuration
            self.key_exists(validation_definition.parameters, "api_url")
            self.key_exists(validation_definition.parameters,
                            "organization_token")
            # call to TF to get entitlements
            entitlements = TFClient.organization_entitlements(
                validation_definition)

            # if they don't have an audit logging key or it's not true, raise
            if not entitlements["audit-logging"]:
                raise ValidationException(
                    "Organization not entitled to audit logging.")

            return True

        except ValidationException as e:
            logging.error(e)
            ErrorFormatter.print_error(str(e))
            sys.exit(1)
        except Exception as e:
            logging.error(e)
            ErrorFormatter.print_error("Error validating arguments")
            sys.exit(1)

    @classmethod
    def key_exists(self, config: Dict[str, Any], key: str) -> None:
        """Raises an error if the key does not exist in the dict

        Args:
            config: dict hodling configuration
            key: key to find
        """
        if key not in config:
            raise ValidationException(
                "Invalid configuration received: key '%s' is missing." % key)

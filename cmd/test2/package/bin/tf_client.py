
import urllib.parse
import configparser
import logging
import requests
from splunklib.modularinput import InputDefinition, ValidationDefinition  # type: ignore
from typing import Any, Optional, Union
from encryptor import Encryptor


class TFClient:
    MASK = "<MASKED>"

    @classmethod
    def audit_trails(self, input_definition: InputDefinition, since: Optional[str] = None, page: int = 0) -> Any:
        """This calls to the audit-trails endpoint.  It will handle pagination and since queries

        Args:
            input_definition: config from Splunk
            since: optional querystring param for last timestamp
            page: page to get

        Returns:
            response from api
        """
        query = {}

        # if since is passed, add to the querystring
        if since is not None:
            query["since"] = since

        # if page is passed, add to the querystring
        if page > 0:
            query["page"] = str(page)

        url = "/api/v2/organization/audit-trail"

        if query != {}:
            # make the query string
            url = url + "?" + urllib.parse.urlencode(query)

        return self.__make_request(input_definition, url)

    @classmethod
    def organization_entitlements(self, config: ValidationDefinition) -> Any:
        """This calls to the organization entitlements endpoint

        Args:
            config: configuration from Splunk

        Returns:
            entitlements dict
        """
        response = self.__make_request(
            config, "/api/v2/organizations?include=entitlement_set")

        # get the included entitlement_set if it is present
        for included in response.get("included", []):
            if included["type"] == "entitlement-sets":
                return included.get("attributes", {})

        return {}

    @classmethod
    def __load_version(self):
        """Attempts to load version from default app.conf
        """
        config = configparser.ConfigParser()
        config.read("../default/app.conf")
        version = "terraform_cloud"
        if "launcher" in config:
            launcher = config["launcher"]
            if "version" in launcher:
                version = version + "/" + launcher["version"]
        return version

    @ classmethod
    def __make_request(self, config: Union[InputDefinition, ValidationDefinition], path: str) -> Any:
        """Makes actual request to API

        Args:
            config: configuration from Splunk
            path: path to request

        Returns:
            object parsed from JSON response
        """
        # config can only be passed in by Splunk
        if not isinstance(config, ValidationDefinition) and not isinstance(config, InputDefinition):
            raise Exception("Invalid configuration")

        name = config.metadata.get("name", "unknown_name")
        server_host = config.metadata.get("server_host", "unknown_host")
        session_key = config.metadata.get("session_key", "unknown_session_key")
        encrypt_token = False

        # load token and api_base
        if isinstance(config, ValidationDefinition):
            # token and api_base are in `parameters` for ValidationDefinition
            token = config.parameters["organization_token"]
            api_url = config.parameters["api_url"]

            # for ValidationDefinitions, we don't want to actually encrypt the token
            # because it could be incorrect
            encrypt_token = False
        elif isinstance(config, InputDefinition):
            # token and api_base are in `inputs` for InputDefinition
            # gets the first stanza because there can evidently be more than one input stanza,
            # but we only care about the first
            name = list(config.inputs.keys())[0]
            name_input = config.inputs[name]
            api_url = name_input["api_url"]
            token = name_input["organization_token"]
            # at this point, we've validated the org token so if it's not masked, we will want to mask it
            encrypt_token = True

        api_base = api_url.strip("/")
        decrypted_token = token

        if token != Encryptor.MASK:
            if encrypt_token:
                Encryptor.encrypt_token(name, token, session_key)
                Encryptor.mask_token(name, api_url, session_key)
        else:
            decrypted_token = Encryptor.get_token(name, session_key)

        # load up the url
        url = urllib.parse.urlparse(api_base + path)

        s = requests.Session()

        s.headers["user-agent"] = self.__load_version()
        s.headers["Authorization"] = "Bearer " + decrypted_token
        s.headers["X-Splunk-Host"] = server_host
        s.headers["X-Splunk-Input-Name"] = name

        r = s.get(url.geturl())

        if r.status_code == 200:
            return r.json()
        elif r.status_code == 401:
            # We got a 401 because their org token is bad
            logging.error("Error 401: Cannot connect to " + url.geturl())
            raise Exception("Error connecting - check organization token")
        else:
            # We got some other status and we are just gonna log it and tell them
            logging.error("Error " + str(r.status_code) +
                          ": Cannot connect to " + url.geturl())
            raise Exception("Error connecting to Terraform Cloud")

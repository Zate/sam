import splunklib.client as client  # type: ignore
import logging


class Encryptor:
    """This encrypts and masks the token using the storage password interface within splunk

    MASK is what we use to mask the token
    """
    MASK = "<MASKED>"
    APP = "terraform_cloud"

    @classmethod
    def encrypt_token(self, name: str, token: str, session_key: str):
        print("FUCK")
        """This method calls to the storage passwords interface and encrypts the organization token

        Args:
            name: name of the input stanza (terraform_cloud://something)
            token: organization token
            session_key: api key to speak to internal Splunk API

        Returns:
            response from api
        """
        if "://" in name:
            _kind, input_name = name.split("://")
        else:
            input_name = name

        args = {"token": session_key}

        service = client.connect(**args)

        try:
            for storage_password in service.storage_passwords:
                if storage_password.username == input_name:
                    # Found existing credentials - we have to delete these or will get error
                    service.storage_passwords.delete(username=storage_password.username)

            service.storage_passwords.create(token, input_name)
        except Exception as e:
            raise Exception(
                "An error occurred updating credentials - check your account permissions. Details: %s" % str(e))

    @classmethod
    def mask_token(self, name: str, api_url: str, session_key: str):
        """This method calls to the inputs interface and masks the organization token

        Args:
            name: name of the input stanza (terraform_cloud://something)
            api_url: terraform api
            session_key: api key to speak to internal Splunk API

        Returns:
            response from api
        """
        if "://" in name:
            kind, input_name = name.split("://")
        else:
            raise Exception("No kind passed to mask_token " + name)

        args = {"token": session_key}
        service = client.connect(**args)

        item = service.inputs.__getitem__((input_name, kind))

        kwargs = {
            "api_url": api_url,
            "organization_token": self.MASK
        }
        try:
            item.update(**kwargs)
        except Exception as e:
            logging.error(e)
            raise Exception("Error updating inputs.conf: %s" % str(e))

    @classmethod
    def get_token(self, name: str, session_key: str):
        """This method calls to the storage passwords interface and gets the stored org token

        Args:
            name: name of the input stanza (terraform_cloud://something)
            session_key: api key to speak to internal Splunk API

        Returns:
            response from api
        """
        if "://" in name:
            kind, input_name = name.split("://")
        else:
            input_name = name

        args = {"token": session_key}

        service = client.connect(**args)

        for storage_password in service.storage_passwords:
            if storage_password.username == input_name:
                # this matches our stanza
                return storage_password.content.clear_password

from splunklib.modularinput import Event  # type: ignore
from datetime import datetime, timezone
import json
from json import JSONEncoder
from typing import Any, Dict


class AuditTrailRequest:
    """Class to hold request data, if available

    Attributes:
        id (str): unique id for the request if present
    """
    id: str

    def __init__(self, request: Dict[str, Any]) -> None:
        self.id = request.get("id", None)

    """Checks equality

    Args:
        other: comparison object

    Returns:
        True for equal, False for not equal
    """

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, AuditTrailRequest):
            return NotImplemented
        return self.id == other.id


class AuditTrailAuth:
    """Class to hold auth data, if available

    Attributes:
        impersonator_id (str): external_id of the admin user impersonating if present
        type (str): type of authentication (Client, Impersonated, System)
        accessor_id (str): external_id of the authenticated user if present
        description (str): username of the authenticated user if present
        organization_id (str): external_id of the organizaiton if present
    """

    def __init__(self, auth: Dict[str, Any]) -> None:
        self.impersonator_id = auth.get("impersonator_id", None)
        self.type = auth.get("type", None)
        self.accessor_id = auth.get("accessor_id", None)
        self.description = auth.get("description", None)
        self.organization_id = auth.get("organization_id", None)

    """Checks equality

    Args:
        other: comparison object

    Returns:
        True for equal, False for not equal
    """

    def __eq__(self, other: object) -> Any:
        if not isinstance(other, AuditTrailAuth):
            return NotImplemented
        return (
            self.impersonator_id,
            self.type,
            self.accessor_id,
            self.organization_id
        ) == (
            other.impersonator_id,
            other.type,
            other.accessor_id,
            other.organization_id
        )


class AuditTrailResource:
    """Class to hold resource data, if available

    Attributes:
        action (str): action taken on resource if present
        meta (dict): key value pairs with metadata if present
        type (str): class name of the resource if present
        id (str): external_id of the resource if present
    """

    def __init__(self, resource: Dict[str, Any]) -> None:
        self.action = resource.get("action", None)
        self.meta = resource.get("meta", {})
        self.type = resource.get("type", None)
        self.id = resource.get("id", None)

    """Checks equality

    Args:
        other: comparison object

    Returns:
        True for equal, False for not equal
    """

    def __eq__(self, other: object) -> Any:
        if not isinstance(other, AuditTrailResource):
            return NotImplemented
        return (self.id == other.id and self.action == other.action)


class AuditTrail:
    """Class to hold audit trail data

    Attributes:
        resource (AuditTrailResource): details about the resource
        request (AuditTrailRequest): details about the request
        auth (AuditTrailAuth): details about the authenticated actor
        timestamp (str): UTC ISO8601 date
        parsed_date (datetime): Python datetime in utc timezone
        version (str): Version of the audit trail schema
        type (str): type of the audit trail
        id (str): unique id of the audit trail
    """

    def __init__(self, params: Dict[str, Any]) -> None:
        self.resource = AuditTrailResource(params.get("resource", {}))
        self.request = AuditTrailRequest(params.get("request", {}))
        self.auth = AuditTrailAuth(params.get("auth", {}))
        self.timestamp = params["timestamp"]
        self.parsed_date = datetime.strptime(self.timestamp, "%Y-%m-%dT%H:%M:%S.000Z").replace(tzinfo=timezone.utc)
        self.version = params.get("version", "0")
        self.type = params.get("type", "Resource")
        self.id = params.get("id", None)

    def to_event(self, stanza: str, host: str) -> Event:
        """Converts an AuditTrail into an modularinput.Event object

        Args:
            stanza: name of the input
            host: the api_url

        Returns:
            event object
        """
        time = "%.3f" % self.parsed_date.timestamp()
        return Event(
            data=json.dumps(self, indent=4, cls=AuditTrailEncoder),
            stanza=stanza,
            time=time,
            host=host,
            source="terraform_cloud",
            sourcetype="terraform_cloud",
            done=True,
            unbroken=True
        )


class AuditTrailEncoder(JSONEncoder):
    """Custom encoder to get json"""

    def default(self, audit_trail: AuditTrail) -> Dict[str, Any]:  # pylint: disable=E0202
        """Removes parsed_date from the dictionary for json

        Args:
            audit_trail: object to encode
        """
        item = audit_trail.__dict__
        if "parsed_date" in item:
            del item["parsed_date"]

        return item

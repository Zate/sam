from pathlib import Path
from datetime import datetime, timezone
import json
import logging
import typing


class Checkpoint:
    """Checkpoint holds the information about the passed in checkpoint file from Splunk

    Attributes:
        path (Path): path to checkpoint directory
        checkpoint_file (Path): path to checkpoint JSON file
        last_run (datetime): datetime of last run for querystring
    """

    def __init__(self, checkpoint_path: str, stanza: str) -> None:
        """This takes a path from Splunk for the checkpoint_dir and tries to see if there is anything
        there yet.  We are specifically looking for a checkpoint.json file.

        That file has a JSON object within it with one variable: last_run  which has the UTC
        timestamp of the last AuditTrail item we parsed

        Args:
            checkpoint_path: posix path to checkpoint directory from Splunk
            stanza: string of format terraform_cloud://STANZA
        """
        self.stanza = stanza.split("://")[1]
        logging.debug("Checkpoint: stanza is " + str(self.stanza))
        self.path = Path(checkpoint_path)
        self.checkpoint_file = self.path / f"{self.stanza}_checkpoint.json"
        logging.debug("Checkpoint: checkpoint_file is at " + str(self.checkpoint_file))
        self.last_run = None

        if self.checkpoint_file.exists():
            logging.debug("Checkpoint: checkpoint_file exists")
            # load up that file there
            with self.checkpoint_file.open() as file:
                details = json.loads(file.read())
                self.last_run = details.get("last_run", datetime.now(timezone.utc).timestamp())
                logging.debug("Checkpoint: last_run at " + str(self.last_run))

    def last_checkpoint(self) -> typing.Union[None, str]:
        """Returns the last_run value in a properly formatted ISO string

        Returns:
            ISO formatted string or None
        """
        if self.last_run is not None:
            dt = datetime.fromtimestamp(self.last_run, timezone.utc)
            return dt.strftime("%Y-%m-%dT%H:%M:%S.000Z")
        return None

    def mark_checkpoint(self, checkpoint: datetime = datetime.now(timezone.utc)) -> None:
        """Creates/updates the file with the last timestamp from audit trail request

        Args:
            checkpoint: last timestamp from audit trail request
        """
        with open(self.checkpoint_file, "w") as file:
            file.write(json.dumps({"last_run": checkpoint.timestamp()}))

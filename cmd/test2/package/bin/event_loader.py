from tf_client import TFClient
from audit_trail import AuditTrail
from checkpoint import Checkpoint
from splunklib.modularinput import InputDefinition  # type: ignore
from typing import List, Callable


class EventLoader:
    @classmethod
    def load_events(cls, input_definition: InputDefinition, log: Callable, stanza: str) -> List[AuditTrail]:
        """Loads events from the API using the checkpoint and config passed in

        Args:
            input_definition: config passed in from Splunk

        Returns:
            list of AuditTrails

        """
        results: List[AuditTrail] = []

        # create the checkpoint object
        checkpoint = Checkpoint(input_definition.metadata["checkpoint_dir"], stanza)
        log("DEBUG", "Checkpoint.last_checkpoint: " + str(checkpoint.last_checkpoint()))

        # start the recursion
        cls.__call_tf_client(input_definition, checkpoint, 1, results)

        log("DEBUG", "Loaded " + str(len(results)) + " events")

        if len(results) > 0:
            # sort the results by reverse timestamp
            results.sort(reverse=True, key=lambda trail: trail.timestamp)
            checkpoint.mark_checkpoint(results[0].parsed_date)

        return results

    @classmethod
    def __call_tf_client(cls,
                         input_definition: InputDefinition,
                         checkpoint: Checkpoint,
                         page: int,
                         results: List[AuditTrail]
                         ) -> None:
        """Loads up the audit trails with the provided configuration for the page

        # checkopoint.last_checkpoint can either be a date or None
        #   either way, we send it to the audit_trails method

        Args:
            input_definition: configuration from Splunk
            checkpoint: checkpoint object
            page: requested page
            results: list to append data to

        """
        # Get the response from the TFClient
        # A response has the structure:
        #
        # {
        #   data: [AuditTrail, AuditTrail....],
        #   pagination: {
        #       next_page: something,
        #   }
        # }
        response = TFClient.audit_trails(
            input_definition,
            since=checkpoint.last_checkpoint(),
            page=page
        )

        # Loop through all of objects in data stanza and add to results
        for trail in response.get("data", []):
            if trail != {}:
                results.append(AuditTrail(trail))

        # Load up the pagination stanza
        pagination = response.get("pagination")

        # Get the next page variable from the stanza
        if pagination["next_page"] is not None and pagination["next_page"] > page:
            cls.__call_tf_client(input_definition, checkpoint, pagination["next_page"], results)

        # let's get out of here
        return

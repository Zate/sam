import xml.sax.saxutils


class ErrorFormatter:
    @classmethod
    def print_error(self, error: str) -> None:
        """Outputs error message as xml to make Splunk happy

        Args:
            error: message to write
        """
        print("<error><message>%s</message></error>" %
              xml.sax.saxutils.escape(error))

import os
from elastalert.alerts import Alerter, BasicMatchString
from notifications_python_client.notifications import NotificationsAPIClient


class GovNotifyAlerter(Alerter):

    required_options = set(['log_file_path', 'email'])

    def __init__(self, rule):
        Alerter.__init__(self, rule)
        self.template_id = os.environ['GOVUK_NOTIFY_TEMPLATE_ID']
        self.email_addresses = os.environ['NOTIFICATION_EMAILS'].split(',')
        api_key = os.environ['GOVUK_NOTIFY_API_KEY']
        self.notifications_client = NotificationsAPIClient(api_key)

    @staticmethod
    def _generate_personalisation(match_items):
        personalisation = {}
        for i, v in enumerate(match_items):
            if v[0] == 'Message':
                personalisation['Message'] = v[1]
            elif v[0] == 'Timestamp':
                personalisation['Timestamp'] = v[1]
            elif v[0] == '_index':
                personalisation['ElasticsearchIndex'] = v[1]
            elif v[0] == '_id':
                personalisation['ElasticsearchId'] = v[1]
            elif v[0] == 'Data':
                personalisation['Filename'] = v[1]['filename']
                personalisation['Reason'] = v[1]['reason']
                personalisation['Organisation'] = v[1]['organisation']
                personalisation['Repo'] = v[1]['repo']
                personalisation['URL'] = v[1]['url']
        return personalisation

    def _send_notification(self, email_address, personalisation):
        return self.notifications_client.send_email_notification(
            email_address=email_address,
            template_id=self.template_id,
            personalisation=personalisation,
            reference=None
        )

    def alert(self, matches):
        # Matches is a list of match dictionaries.
        # It contains more than one match when the alert has
        # the aggregation option set
        for match in matches:
            personalisation = self._generate_personalisation(match.items())
            for email_address in self.email_addresses:
                self._send_notification(
                    email_address, personalisation)

            with open(self.rule['log_file_path'], 'a') as output_file:
                # basic_match_string will transform the match into the default
                # human readable string format
                # https://github.com/Yelp/elastalert/blob/3931d7feaf0d07b6531fb53042b9284bb46712ce/elastalert/alerts.py#L128
                match_string = str(BasicMatchString(self.rule, match))
                output_file.write(match_string)

    # get_info is called after an alert is sent to get
    # data that is written back to Elasticsearch in the field "alert_info"
    # It should return a dict of information relevant to what the alert does
    def get_info(self):
        return {'type': 'GovUK Notify Alerter',
                'email': self.rule['email'],
                'log_file_path': self.rule['log_file_path']}

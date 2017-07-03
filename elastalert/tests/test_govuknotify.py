import pytest
from modules.govuknotify import GovNotifyAlerter


@pytest.mark.parametrize(('match_items', 'expected'), [
    (
        [
            ('Message', 'some message'),
            ('Timestamp', '12:30'),
            ('_index', 'githubintegration'),
            ('_id', '12345'),
            ('Data', {
                'filename': 'secret.txt',
                'reason': 'it is a secret',
                'organisation': 'homeoffice',
                'repo': 'greatrepo',
                'url': 'https://github.com/homeoffice/greatrepo',
            }),
        ],
        {
            'Message': 'some message',
            'Timestamp': '12:30',
            'ElasticsearchIndex': 'githubintegration',
            'ElasticsearchId': '12345',
            'Filename': 'secret.txt',
            'Reason': 'it is a secret',
            'Organisation': 'homeoffice',
            'Repo': 'greatrepo',
            'URL': 'https://github.com/homeoffice/greatrepo',
        }
    ),
])
def test_personalisation(match_items, expected):
    personalisation = GovNotifyAlerter._generate_personalisation(
        match_items=match_items)
    assert personalisation == expected

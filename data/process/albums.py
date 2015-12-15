from mrjob.job import MRJob
from mrjob.step import MRStep
import json

class MRDeriveAlbumDataset(MRJob):

    def mapper(self, _, line):
        data = json.loads(line)
        page = data["page"]
        for album in data["albums"]:
            yield page, album

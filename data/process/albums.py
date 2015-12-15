import json
import os.path
import re

import mrjob.protocol
from mrjob.job import MRJob
from mrjob.step import MRStep

IMAGE_ID_RE = re.compile(r"[0-9]+")


class MRDeriveAlbumDataset(MRJob):
    OUTPUT_PROTOCOL = mrjob.protocol.JSONValueProtocol

    def mapper(self, _, line):
        data = json.loads(line)
        ## yield each album after processing
        for album in data["albums"]:
            album["album_id"] = int(os.path.basename(album["URI"]))

            ## only include image id's and not full info
            images = []
            for image in album["images"]:
                images.append(int(IMAGE_ID_RE.findall(os.path.basename(image["uri"]))[0]))
            album["images"] = images
            yield album["album_id"], album

    def reducer(self, _, album):
        for val in album:
            yield None, val

if __name__ == '__main__':
    MRDeriveAlbumDataset.run()

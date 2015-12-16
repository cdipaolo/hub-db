import json
import os.path
import re

import mrjob.protocol
from mrjob.job import MRJob

IMAGE_ID_RE = re.compile(r"[0-9]+")


class MRDeriveAlbumDataset(MRJob):
    """MRDeriveAlbumDataset
    this MRJob MapReduce job takes the raw_data
    from crawling pornhub and generates an album
    specific dataset fromt that data.
    """
    OUTPUT_PROTOCOL = mrjob.protocol.JSONValueProtocol

    def mapper(self, _, line):
        """mapper
        the map step takes each album from the page,
        finds it's album_id from the URL scheme and
        saves it, and maps each image to it's image
        id instead of carrying around all the image
        metadata (that's kept in a separate dataset)
        """
        data = json.loads(line)
        ## yield each album after processing
        for album in data["albums"]:
            album["album_id"] = int(os.path.basename(album["URI"]))

            ## only include image id's and not full info
            images = []
            for image in album["images"]:
                images.append(int(IMAGE_ID_RE.findall(os.path.basename(image["uri"]))[0]))
            album["images"] = images
            yield None, album

    def reducer(self, _, album):
        """reducer
        the reduce step just passes all map
        output to the OUTPUT as a json object
        """
        for val in album:
            yield None, val

if __name__ == '__main__':
    MRDeriveAlbumDataset.run()

import json
import os.path
import re

import mrjob.protocol
from mrjob.job import MRJob
from mrjob.step import MRStep

IMAGE_ID_RE = re.compile(r"[0-9]+")


class MRDeriveImageDataset(MRJob):
    """MRDeriveImageDataset
    this MRJob MapReduce class takes in the raw_data
    from crawling pornhub and generates an image specific
    dataset while still holding onto information about
    which albums the images were a part of
    """
    OUTPUT_PROTOCOL = mrjob.protocol.JSONValueProtocol

    def mapper(self, _, line):
        """mapper
        the map step splits images from their albums,
        storing the segment, album_id, and album_title
        for use later and yielding the image by itself
        with no key
        """
        data = json.loads(line)
        ## yield each image after processing
        for album in data["albums"]:
            album_id = int(os.path.basename(album["URI"]))

            ## yield images from the albums
            for image in album["images"]:
                image["segment"] = album["segment"]
                image["album_id"] = album_id
                image["album_title"] = album["title"]
                image["image_id"] = int(IMAGE_ID_RE.findall(os.path.basename(image["uri"]))[0])
                yield None, image

    def reducer(self, _, image):
        """reducer
        the reduce step just passes all map
        output to the OUTPUT as a json object
        """
        for val in image:
            yield None, val


if __name__ == "__main__":
    MRDeriveImageDataset.run()

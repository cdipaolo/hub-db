import json

import mrjob.protocol
from mrjob.job import MRJob

class MRDeriveCommentsDataset(MRJob):
    """MRDeriveCommentsDataset
    this MRJob MapReduce class takes in the images dataset
    and creates a dataset of comments while still holding
    onto the segment, image_id, album_id, and other metadata
    of the images

    note that only the popular_comments were crawled, so this
    dataset does not contain all pornhub comments on images ever
    """
    OUTPUT_PROTOCOL = mrjob.protocol.JSONValueProtocol

    def mapper(self, _, line):
        """mapper
        the map step takes in an image json
        object and maps it to comment information,
        holding onto information about the original
        image
        """
        image = json.loads(line)
        ## yield each comment
        if "popular_comments" not in image:
            return
        for comment in image["popular_comments"]:
            comment["segment"] = image["segment"]
            comment["image_id"] = image["image_id"]
            comment["album_id"] = image["album_id"]
            comment["album_title"] = image["album_title"]
            yield comment["image_id"], comment

    def reducer(self, _, comment):
        """reducer
        the reduce step just passes all map
        output to the OUTPUT as a json object
        """
        for val in comment:
            yield None, val

if __name__ == "__main__":
    MRDeriveCommentsDataset.run()

import json

import mrjob.protocol
from mrjob.job import MRJob

class MRDeriveTagsFrequencies(MRJob):
    """MRDeriveTagsFrequencies
    this MRJob MapReduce job takes the albums
    dataset and generates the frequencies
    of all tags

    this MR job doesn't retain tags with
    lengths less than 3 or counts of less
    than 6
    """

    def mapper(self, _, line):
        """mapper
        the map step takes in an album and
        maps it to a reduced format comprised
        of the tag and all tags in the same
        album including the original tag
        """
        album = json.loads(line)
        ## record associated tags
        for tag in album["tags"]:
            if len(tag) < 3:
                return
            yield tag.lower(), 1

    def reducer(self, key, counts):
        """reducer
        the reduce step takes all the tags
        and returns a key (tag) and the
        number of times the tag has been
        seen
        """
        s = sum(counts)
        if s < 6:
            return
        yield key, s


if __name__ == "__main__":
    MRDeriveTagsFrequencies.run()

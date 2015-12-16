import json

import mrjob.protocol
from mrjob.job import MRJob

class MRDeriveTagsDataset(MRJob):
    """MRDeriveTagsDataset
    this MRJob MapReduce job takes the albums
    dataset and generates associations between
    which tags are seen with which other tags
    in terms of garph links. Weights are
    given as the frequency of occurences

    this MR job doesn't retain tags with
    lengths less than 3
    """

    def mapper(self, _, line):
        """mapper
        the map step takes in an album and
        maps it to a reduced format comprised
        of the tag and all tags in the same
        album including the original tag
        """
        album = json.loads(line)
        ## filter out the small tags
        album["tags"] = [tag for tag in album["tags"] if len(tag) > 2]
        ## record associated tags
        for i, tag in enumerate(album["tags"]):
            for other in album["tags"]:
                yield tag.lower(), other.lower()

    def reducer(self, key, tags):
        """reducer
        the reduce step takes all the tags
        and returns a key (tag) and an array
        of 2-tuples of tags seen with the key
        as well as the number of times it has
        been seen

        note that the tag list contains
        the key, and that count corresponds to
        the total number of times the tag has
        been seen
        """
        seen = {}
        for tag in tags:
            if tag in seen:
                seen[tag] += 1
            else:
                seen[tag] = 1
        ## now generate a list from the dictionary
        l = []
        for tag, count in seen.iteritems():
            l.append((tag, count))
        yield key, l


if __name__ == "__main__":
    MRDeriveTagsDataset.run()

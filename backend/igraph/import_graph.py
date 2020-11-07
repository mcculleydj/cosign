import argparse
import os
import pymongo
import igraph
import dotenv


def main():
    g = igraph.Graph()
    g.add_vertices(members_collection.count_documents({}))

    for cell in adjacency_cells_collection.find():
        s = int(cell['position'].split('_')[0]) - 1
        t = int(cell['position'].split('_')[1]) - 1
        g.add_edge(s, t, weight=len(cell['billNumbers']))

    print(g)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-d', dest='is_dev', action='store_true')
    args = parser.parse_args()

    dotenv.load_dotenv()

    if args.is_dev:
        uri = 'mongodb://localhost:27017'
    else:
        uri = f'mongodb+srv://{os.getenv("MONGO_USER")}:{os.getenv("MONGO_PASSWORD")}@cluster0.wht7g.mongodb.net/admin?retryWrites=true&w=majority'

    client = pymongo.MongoClient(uri)
    members_collection = client.cosign.members
    adjacency_cells_collection = client.cosign.cells

    main()

    client.close()

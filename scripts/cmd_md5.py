import sys, json, hashlib

args = sys.argv[1:]
if len(args) == 0: sys.exit("not input args")

d = {"commandline": args}
bts = json.dumps(d).replace(" ", "").encode('utf8')
# print(d)
# print(bts)

hash_object = hashlib.md5(bts)
print(hash_object.hexdigest())

import os

files = os.listdir("./")
# print(files)
for f in files:
    if f.endswith("_grpc.pb.go"):
        openedFile = open(f, "r")
        fileContent = openedFile.read()
        truncatedFileContent = fileContent.split("Server is the server API")[0]
        openedFile.close()
        with open(f, "w") as toWrite:
            toWrite.write(truncatedFileContent)
Write path:

1. Client → CreateFile(file_path)

2. Compute chunk_index from offset

3. Client → GetChunkLocations(file_path, chunk_index)

4. IF CHUNK_NOT_FOUND:
       Client → AllocateChunk(file_path)

5. Client → WriteChunk(chunk_id, offset, data)
       → send to all chunkservers

6. Collect responses → return success/failure

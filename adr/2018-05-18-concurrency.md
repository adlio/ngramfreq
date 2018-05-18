# Concurrency Architecture Decision Record

**Status**: Proposed

**Context**: When processing multiple files, especially from different physical disks, running Scanner.Scan concurrently should improve performance.

**Decision**: Evaluate options. The existing design, using a non thread-safe map is ~2X faster when run on a single file than a version which used `sync.Map`. We suspect this is due to locking overhead (`sync.Map` relies on a mutex). It's possible that sufficiently parallel operations the locking overhead would be overcome by the increased parallelized performance. It would be useful to try other concurrent approaches, such as concurrent n-gram tokens being sent to a channel with a single consumer/scorer.

**Consequences**: Performance is optimized for single machine, single-disk environments, but may underperform when being run with large numbers of files from many disks.


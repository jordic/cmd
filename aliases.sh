#!/bin/bash

function large_files {
    find . -type f -print0 | xargs -0 du -h | sort -hr | head -40
}

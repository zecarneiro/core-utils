import sys

from vendor.file_utils import file_exist


def main():
    args = sys.argv[1:]
    if len(args) == 1:
        file_exist()

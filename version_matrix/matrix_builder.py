import datetime
import threading
import logging

from . import release_versions
from . import docker
from . import constant


def build_matrix(username: str, password: str, repository: str) -> dict:
    """
    Build the version matrix by spawning checks for validity
    of the version detected against previous builds, their age
    and if newer patch versions are available

    :param username:
    :param password:
    :param repository:
    :return: None
    """

    threads = []
    matrix = {"version": [], "include": []}
    failures = []

    for version in release_versions.list_all_versions():
        logging.info("Spawning check for version %s", version)

        thread = threading.Thread(
            target=_check_version,
            args=(version, matrix, failures, username, password, repository),
            name="BuildThread-" + version
        )

        threads.append(thread)
        thread.start()

    logging.info("Waiting for all threads to finish...")
    for thread in threads:
        thread.join()

    return matrix


def _check_version(version_number: str, matrix: dict, failures: list, username: str, password: str, repository: str):
    """
    Check a particular version for freshness

    :param version_number:
    :param matrix:
    :param failures:
    :param :username:
    :param password:
    :param repository:
    :return:
    """

    minimum_age = datetime.datetime.now() - datetime.timedelta(days=constant.MAX_AGE_IN_DAYS)
    version_number_zts = version_number + "-zts"

    try:
        version_metadata = release_versions.fetch_version_metadata(version_number)
        logging.debug(version_metadata)

        metadata = docker.fetch_metadata(username, password, repository, version_number)
        logging.debug(metadata)

        for platform in metadata:
            last_modified = metadata[platform]

            if last_modified < version_metadata["release_date"] or last_modified < minimum_age:
                logging.info("Appending %s to build list", version_number)

                matrix["version"].append(version_number)
                matrix["include"].append({
                    "version": version_number,
                    "asset_url": version_metadata["release_asset_url"],
                })

                matrix["version"].append(version_number_zts)
                matrix["include"].append({
                    "version": version_number_zts,
                    "asset_url": version_metadata["release_asset_url"],
                })

                return

    except Exception:
        failures.append(version_number)

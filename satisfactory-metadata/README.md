# An helper that sync Satisfactory meta to a Postges database

This is an helper that sync Satisfactory meta to a Postges database, in order to be able to query it from Grafana.

This is an helper for Satisfactory that use [Ficsit Remote Montioring mod](https://ficsit.app/mod/B9bEiZFtaaQZHU) and write the data to a Postgres database.

This work is a fork of the work of [Jeff Wong](https://github.com/featheredtoast/) on the [Ficsit Remote Monitoring Companion Bundle](https://github.com/featheredtoast/satisfactory-monitoring) that I have adapted to my usecase.

Basicaly I have keep the part that write to the database; and the Dashboards.
I have edit the database sync part to take advantage of the kubernetes jobs.

And added a wrapper for building this as a Docker image.

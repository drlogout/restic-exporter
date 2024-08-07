# Restic exporter

Small prometheus exporter that does a single thing: It exports the time at which
the last [restic](github.com/restic/restic) backup has been done. That's it.

## Usage

This is not a standalone exporter - rather it generates `.prom` files that can
be picked up by [`node-exporter`](github.com/prometheus/node_exporter).

In order to use it, first setup `node-exporter` to watch a directory. If you are
on debian based systems, this can ususally be done by adding the
`-collector.textfile.directory` option to `/etc/default/prometheus-node-exporter`:

```txt
ARGS="-collector.textfile.directory=\"/var/lib/prometheus/node-exporter/\""
```

Next, create a .env file for `restic-exporter` with RESTIC_REPOSITORY, RESTIC_PASSWORD and additional variables needed:

```sh
# .env
RESTIC_REPOSITORY=
RESTIC_PASSWORD=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=
```

Now, you can let `restic-exporter` generate a `.prom` file which in turn is picked
up and exposed to prometheus by `node-exporter`.

```bash
restic-exporter -name arbitrary-name-here
```

This will generate a single stat:

```txt
restic_snapshot_timestamp{name="arbitrary-name-here"} 1.599849001e+09
```

## Why not make this a 'real' exporter

I went the `node-exporter` route here because I personally store backups on a NAS
with spinning hard drives that are spun down most of the time. Thus, I needed a way
to check the backup status at a specific time where the drives are already spun up.

The easiest way to achieve this was the file-generation route in combination with
a cronjob.

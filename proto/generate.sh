cd proto
buf generate
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/noble-assets/paramauthority/* ./
rm -rf github.com
rm -rf noble

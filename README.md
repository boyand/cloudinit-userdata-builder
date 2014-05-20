cloudinit-userdata-builder
====================

### About 

This script bundles a number of files into cloudinit supportable userdata, which can be passed to your instances (which may be running in your local DC, on AWS or on other cloud providers)

### Installation

    go get github.com/boyand/cloudinit-userdata-builder

### Usage

This binary supports two flags by default: config file path and encode flag for output formatting. In order to use it, your need to supply a config file for your userdata using the following format: 

    cat userdata.conf
```
{
    "cloud_init_parts": [
        [
            "/path/to/part-handler.py",
            "text/part-handler"
        ],
        [
            "/path/to/cloudinit.txt",
            "text/cloud-config"
        ],
        [
            "/path/to/your/script.sh",
            "text/cloud-boothook"
        ]
    ]
}
```

If you want to use the binary to build drop-in bundle for AWS, you can run it with: 

    cloudinit-userdata-builder --config="/path/to/userdata.conf" --encode

 



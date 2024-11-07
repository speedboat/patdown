# patdown

> Remotely predicts and identifies the presence of EDR/XDR solutions on networks

<p  align="center">
    <img  src="https://i.imgur.com/ynSGB7L.png"  width="500"  title="hover text">
</p>

## Abstract
patdown is an EDR/XDR fingerprinting utility used for remotely predicting defense mechanisms in use on a network.

This allows you to forecast the security posture of a network during the earliest stages of access, or even prior to any access at all.

Fingerprinting is achieved via the probing of DNS servers to determine whether they have resolved domains associated with various EDR/XDR solutions.

**Example**: if a network's resolver has `assets-public.falcon.crowdstrike.com` cached, chances are the *CrowdStrike Falcon* EDR solution is present on the network.

These DNS servers can be specified as arguments (most effective), or patdown can automatically retrieve and analyze the authoritative nameservers of a target with the `-d` flag.

> ⚠️  Authoritative nameservers are rarely used as egress resolvers for networks and are not as reliable for fingerprinting EDR/XDR, making them prone to false positives.

## Installation
Retrieve a binary corresponding to your architecture from **Releases** 

*or*

`git clone https://github.com/speedboat/patdown.git ; cd patdown/cmd/patdown ; go build -o patdown main.go ; ./patdown -h`

## Usage
```
  d | target fqdn (not as reliable, prone to false positives)
  n | nameserver to query (can be specified multiple times)
  v | enable verbosity [false]
  t | threads [5]
  s | delay between requests in milliseconds, per thread [250]

e.g.
    patdown -d target.network
    patdown -n egress.ns.target.network -n another.egress.ns.target.network
    patdown -n dc.target.network -v -t 25
```

## Currently Identified Vendors/Solutions:
- [x] **CrowdStrike** Falcon
- [x] **Microsoft** Defender for Endpoint
- [x] **VMWare** Carbon Black
- [x] **Check Point** Harmony
- [x] **Cybereason** EDR
- [x] **Trellix** EDR
- [x] **Palo Alto Networks** Cortex XDR
- [x] **SentinelOne** Singularity
- [x] **Symantec** Endpoint Security
- [x] **Tanium** EDR
- [x] **Nextron** Aurora
- [x] **Trend Micro** Endpoint Sensor
- [x] **Rapid7** InsightIDR
- [ ] **ESET** Inspect
- [ ] **Harfanglab** EDR
- [ ] **Limacharlie** EDR
- [ ] **Elastic** Security
- [ ] **Qualys** EDR
- [ ] **Uptycs** XDR
- [ ] **WatchGuard** EDR

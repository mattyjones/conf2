# Conf2
Conf2 is an implementation for two emerging standards in the microservice management space:

* YANG (RFC6020) - Interface Definition Language (IDL) covering configuration, metrics, RPCs and events. YANG's strength's include human-readability, developer reusability and the first IDL to cover events over a web service.
* RESTCONF (RFC pending) - A RESTCONF service is just a RESTful service following a few conventions such that consumers of RESTCONF services may be unaware they are using RESTCONF. Two services however with the same YANG will have the exact same RESTCONF API.

Because Conf2 is designed as a library, you can make any running process RESTCONF capable without running any other services.

Benefits of YANG and RESTCONF:
* Standards compliance means automatic infrastrucure integration with other standards based controller systems.
* Receiving configuration through the network reduces reliance file-based configuration tools such as Puppet or Chef
* Exporting health and metrics data through the network reduces reliance on log scraping tools like Splunk
* Sending alerts as they happen to subscribed systems reduces reliance on poll-based systems like watchdog or Nagios.
* Exporting operational functions (e.g. cache-clearing or traffic routing) through the network reduces reliance on tools like Ansible.

Benefits and features of Conf2:
* Written in the Go with C-compatible API enabling support for languages including Java, PHP, Python, Ruby, JavaScript, C, and C++
* Experimental support for Java
* No dependencies beyond Go Standard Library
* API designed to integrate into any existing codebase without modification.
* No code generation
* Access to meta-data for model-driven UIs and tools
* Ability to add custom protocols beyond RESTCONF including NETCONF, SNMP, Weave or proprietary protocols
* Ability to add custom formats beyond JSON or XML
* Experimental support for distributed UI using Web Components

Code Examples:
* [HelloGo](examples/helloGo/README.md) - Basic Go application with RESTCONF API
* [Todo](examples/todo/README.md) - Todo Go application with RESTCONF API

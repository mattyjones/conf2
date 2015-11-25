# Conf2
Conf2 is an implementation for two emerging standards in the microservice management space:

* YANG (RFC6020) - An Interface Definition Language (IDL) covering configuration, metrics, RPCs and events. YANG's strengths include human-readability, developer reusability and support for RPCs and events.
* RESTCONF (RFC pending) - A RESTCONF service is a RESTful service that follows a few conventions. Consumers of RESTCONF services may be unaware they are using RESTCONF. Two services with the same YANG will have the same RESTCONF API.

Conf2 is designed as a library allowing you can make any running process a RESTCONF-capable without running any other services.

Benefits of YANG and RESTCONF:
* Standards compliance means automatic infrastrucure integration with other standards based controller systems.
* Receiving configuration through the network obviates file-based configuration tools such as Puppet or Chef
* Exporting health and metrics data through the network obviates log scraping tools like Splunk
* Sending alerts as they happen to subscribed systems obviates poll-based systems like watchdog or Nagios.
* Exporting operational functions (e.g. cache-clearing or traffic routing) through the network obviates tools like Ansible.

Benefits and features of Conf2:
* Written in the Go with C-compatible API enabling support for languages including Java, PHP, Python, Ruby, JavaScript, C, and C++
* Experimental support for Java
* No dependencies beyond Go Standard Library
* API designed to integrate into any existing codebase without modification.
* No code generation
* Access to meta-data for model-driven UIs and tools
* Ability to add custom protocols beyond RESTCONF including NETCONF, SNMP, Weave or proprietary protocols
* Ability to add custom formats beyond JSON or XML
* Experimental support for distributed UI using Web Components (emerging W3C standard).

Code Examples:
* [HelloGo](examples/helloGo/hello.go) - Basic Go application that says hello to you
* [Todo](examples/todo) - Todo Go application

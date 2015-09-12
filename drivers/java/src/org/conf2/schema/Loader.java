package org.conf2.schema;

/**
 *
 */
public interface Loader {

    public Module loadModule(StreamSource source, String resource);
}

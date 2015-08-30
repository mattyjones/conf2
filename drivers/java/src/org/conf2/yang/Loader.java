package org.conf2.yang;

/**
 *
 */
public interface Loader {

    public Module loadModule(StreamSource source, String resource);
}

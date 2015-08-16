package org.conf2.yang;

/**
 *
 */
public interface Loader {

    public Module loadModule(DataSource source, String resource);
}

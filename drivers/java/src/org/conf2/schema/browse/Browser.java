package org.conf2.schema.browse;

import org.conf2.schema.Module;

public interface Browser {
    public Selection getRootSelector();
    public Module getModule();
}
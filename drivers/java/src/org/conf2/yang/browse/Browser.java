package org.conf2.yang.browse;

import org.conf2.yang.Module;

public interface Browser {
    public Selection getRootSelector();
    public Module getModule();
}
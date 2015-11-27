package org.conf2.schema.browse;

import org.conf2.schema.MetaCollection;
import org.conf2.schema.Module;

public interface Browser {
    public Selection getSelector(BrowsePath path);
    public MetaCollection getSchema();
}
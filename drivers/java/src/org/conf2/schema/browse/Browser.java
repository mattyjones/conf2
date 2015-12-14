package org.conf2.schema.browse;

import org.conf2.schema.MetaCollection;
import org.conf2.schema.Module;

public interface Browser {
    public Node getNode();
    public MetaCollection getSchema();
}
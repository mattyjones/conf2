package org.conf2.schema.browse;

import org.conf2.schema.Meta;
import org.conf2.schema.MetaCollection;

/**
 *
 */
public class WalkState {
    public MetaCollection meta;
    public Meta position;
    public Node node;
    public BrowseValue[] key;
    public boolean insideList;
    public boolean found;
}

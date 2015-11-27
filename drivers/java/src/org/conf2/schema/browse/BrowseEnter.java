package org.conf2.schema.browse;

import org.conf2.schema.MetaCollection;

/**
 *
 */
public interface BrowseEnter {
    public Node Enter(Selection sel, MetaCollection meta, boolean create);
}

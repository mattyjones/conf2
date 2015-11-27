package org.conf2.schema.browse;

/**
 *
 */
public interface BrowseIterate {
    public Node Iterate(Selection sel, BrowseValue[] key, boolean create, boolean isFirst);
}

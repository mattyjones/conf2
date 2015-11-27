package org.conf2.schema.browse;

import org.conf2.schema.Meta;

/**
 *
 */
public interface BrowseRead {
    BrowseValue Read(Selection sel, Meta meta);
}

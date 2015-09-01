package org.conf2.yang.browse;

import org.conf2.yang.Meta;

import java.util.Iterator;

/**
 *
 */
public interface WalkController {
    public boolean selectionIterator(Selection s, int level, boolean isFirst);

    public Iterator<Meta> containerIterator(Selection s, int level);
}

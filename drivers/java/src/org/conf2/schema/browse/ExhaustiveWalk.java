package org.conf2.schema.browse;

import org.conf2.schema.Meta;
import org.conf2.schema.MetaCollectionIterator;

import java.util.Iterator;

/**
 * Walk every node and every leaf given a selection.
 */
public class ExhaustiveWalk implements WalkController {
    private final static String[] NO_KEYS = new String[0];

    public boolean selectionIterator(Selection s, int level, boolean isFirst) {
        return s.Iterate.Iterate(NO_KEYS, isFirst);
    }

    public Iterator<Meta> containerIterator(Selection s, int level) {
        return new MetaCollectionIterator(s.meta);
    }
}

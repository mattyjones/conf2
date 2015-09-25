package org.conf2.schema.browse;

import org.conf2.schema.Meta;
import org.conf2.schema.MetaUtil;

import java.util.Collections;
import java.util.Iterator;

/**
 *
 */
class PathWalkController implements WalkController {
    private static BrowseValue[] NO_KEYS = new BrowseValue[0];
    BrowsePath path;
    Selection target;

    PathWalkController(BrowsePath path) {
        this.path = path;
    }

    public boolean selectionIterator(Selection s, int level, boolean isFirst) {
        if (level == path.segments.length) {
            if (path.segments[level - 1].keys.length == 0) {
                target = s;
                return false;
            }
            if (!isFirst) {
                target = s;
                return false;
            }
        }
        if (isFirst && level > 0 && level <= path.segments.length) {
            return s.Iterate.Iterate(NO_KEYS, isFirst);
        }
        return false;
    }

    public Iterator<Meta> containerIterator(Selection s, int level) {
        if (level >= path.segments.length) {
            target = s;
            return Collections.EMPTY_SET.iterator();
        }

        Meta position = MetaUtil.findByIdent(s.meta, path.segments[level].ident);
        return Collections.singleton(position).iterator();
    }
}
package org.conf2.yang.browse;

import org.conf2.yang.*;

import java.util.Collections;
import java.util.Iterator;

/**
 *
 */
public class Walk {
    public static Selection walk(Selection selection, BrowsePath path) {
        PathWalkController controller = new PathWalkController(path);
        walk(selection, controller);
        return controller.target;
    }


    public static void walk(Selection selection, WalkController controller) {
        Walk.walk(selection, controller, 0);
    }

    static void walk(Selection s, WalkController c, int level) {
        if (s.meta instanceof MetaList && ! s.insideList) {
            MetaList l = (MetaList) s.meta;
            boolean hasMore = c.selectionIterator(s, level, true);
            while (hasMore) {
                s.insideList = true;
                walk(s, c, level);
                hasMore = c.selectionIterator(s, level, true);
            }
        } else {
            Selection child;
            Iterator<Meta> i = c.containerIterator(s, level);
            while (i.hasNext()) {
                s.position = i.next();

                // TODO: resolve choice

                if (MetaUtil.isLeaf(s.position)) {
                    BrowseValue v = new BrowseValue();
                    s.Read.Read(v);
                } else {
                    child = s.Enter.Enter();
                    if (!s.found) {
                        child.meta = (MetaCollection) s.position;
                    }
                    walk(child, c, level + 1);

                    if (s.Exit != null) {
                        s.Exit.Exit();
                    }
                }
            }
        }
    }
}

interface WalkController {
    public boolean selectionIterator(Selection s, int level, boolean isFirst);

    public Iterator<Meta> containerIterator(Selection s, int level);
}

class PathWalkController implements WalkController {
    private static String[] NO_KEYS = new String[0];
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
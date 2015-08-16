package org.conf2.yang.driver;

import org.conf2.yang.browse.BrowseValue;
import org.conf2.yang.browse.EditOperation;
import org.conf2.yang.browse.Selection;
import org.conf2.yang.Choice;
import org.conf2.yang.MetaCollection;
import org.conf2.yang.MetaUtil;

/**
 *
 */
public class BrowserAdaptor {

    static Selection enter(Selection s, String ident) {
        s.position = MetaUtil.findByIdent(s.meta, ident);
        Selection child = s.Enter.Enter();
        if (child.found) {
            child.meta = (MetaCollection) s.position;
        }
        return child;
    }

    static boolean iterate(Selection s, String ident, String encodedKeys, boolean first) {
        s.position = MetaUtil.findByIdent(s.meta, ident);
        String[] keys = encodedKeys.split(" ");
        return s.Iterate.Iterate(keys, first);
    }

    static BrowseValue read(Selection s, String ident) {
        s.position = MetaUtil.findByIdent(s.meta, ident);
        BrowseValue val = new BrowseValue();
        s.Read.Read(val);
        return val;
    }

    static void edit(Selection s, String ident, int opCode, BrowseValue val) {
        s.position = MetaUtil.findByIdent(s.meta, ident);
        EditOperation op = EditOperation.values()[opCode];
        s.Edit.Edit(op, val);
    }

    static void exit(Selection s, String ident) {
        s.position = MetaUtil.findByIdent(s.meta, ident);
        s.Exit.Exit();
    }

    static String choose(Selection s, String choiceIdent) {
        Choice choice = (Choice) MetaUtil.findByIdent(s.meta, choiceIdent);
        s.position = s.Choose.Choose(choice);
        return s.position.getIdent();
    }
}

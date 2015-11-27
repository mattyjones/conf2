package org.conf2.schema.driver;

import org.conf2.schema.*;
import org.conf2.schema.browse.*;
import org.conf2.schema.comm.DataReader;
import org.conf2.schema.comm.DataWriter;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;

/**
 *
 */
public class BrowserAdaptor {
    private final static BrowseValue[] NO_KEYS = new BrowseValue[0];

    public static Selection getSelector(Browser browser, String path) {
        BrowsePath p = new BrowsePath(path);
        return browser.getSelector(p);
    }

    public static Node enter(Selection s, String ident, boolean create) {
        s.state.position = updatePosition(s.state.meta, ident);
        Node child = s.node.Select(s, (MetaCollection) s.state.position, create);
//        if (child != null) {
//            child.state.meta = (MetaCollection) s.state.position;
//        }
        return child;
    }

    public static Node iterate(Selection s, ByteBuffer encodedKeyValues, boolean create, boolean first) {
        BrowseValue[] key;
        MetaList meta = (MetaList)s.state.meta;
        if (encodedKeyValues != null) {
            DataReader r = new DataReader(encodedKeyValues);
            DataType[] keyTypes = meta.getKeyDataTypes();
            key = r.readValues(keyTypes);
        } else {
            key = NO_KEYS;
        }
        return s.node.Next(s, meta, create, key, first);
    }

    public static ByteBuffer read(Selection s, String ident) {
        s.state.position = updatePosition(s.state.meta, ident);
        BrowseValue val = s.node.Read(s, s.state.position);
        ByteBuffer data = null;
        if (val != null) {
            DataWriter w = new DataWriter();
            w.writeValue(val);
            data = w.getByteBuffer();
        }
        return data;
    }

    static Meta updatePosition(MetaCollection meta, String ident) {
        Meta position = null;
        if (ident != null) {
            int choiceSep = ident.indexOf('/');
            if (choiceSep >= 0) {
                String choiceIdent = ident.substring(0, choiceSep);
                Choice c = (Choice) MetaUtil.findByIdent(meta, choiceIdent);
                if (c != null) {
                    int caseSep = ident.indexOf('/', choiceSep + 1);
                    String caseIdent = ident.substring(choiceSep + 1, caseSep);
                    ChoiceCase ccase = c.getCase(caseIdent);
                    if (ccase != null) {
                        position = new MetaCollectionIterator(ccase).next();
                    }
                }
            } else {
                position = MetaUtil.findByIdent(meta, ident);
            }
            if (position == null) {
                throw new DriverError("Could not find meta " + ident + " in " + meta.getIdent());
            }
        }
        return position;
    }

    public static void edit(Selection s, String ident, ByteBuffer encodedValue) {
        s.state.position = updatePosition(s.state.meta, ident);
        BrowseValue val = null;
        if (encodedValue != null) {
            DataReader r = new DataReader(encodedValue);
            val = r.readValue(((HasDataType) s.state.position).getDataType());
        }
        s.node.Write(s, s.state.position, val);
    }

    public static void event(Selection s, int eventId) {
        DataEvent e = DataEvent.values()[eventId];
        s.node.Event(s, e);
    }

    public static String choose(Selection s, String choiceIdent) {
        Choice choice = (Choice) MetaUtil.findByIdent(s.state.meta, choiceIdent);
        s.state.position = s.node.Choose(s, choice);
        return s.state.position.getIdent();
    }
}

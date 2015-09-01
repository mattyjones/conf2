package org.conf2.yang;

/**
 * Construct meta tree from yang-1.0.yang by hand, eventually generate it from template.
 */
public class YangModule {
    public static Module YANG;
    static {
        Module yang = new Module("yang");
        yang.setNamespace("http://yang.org/yang");
        yang.setPrefix("yang");
        // TODO: set revision
        {
            Grouping defHeader = grouping(yang, "def-header");
            leaf(defHeader, "ident", "string");
            leaf(defHeader, "description", "string");
        }
        {
            Grouping typeGrouping = grouping(yang, "type");
            Container type = contr(typeGrouping, "type");
            leaf(type, "ident", "string");
            leaf(type, "range", "string");
            leafList(type, "enumeration", "string");
        }
        {
            Grouping tdGrouping = grouping(yang, "groupings-typedefs");
            MetaList groupings = list(tdGrouping, "groupings");
            groupings.setEncodedKeys("ident");
            groupings.addMeta(new Uses("def-header"));
            groupings.addMeta(new Uses("groupings-typedefs"));
            groupings.addMeta(new Uses("containers-lists-leafs-uses-choice"));

            MetaList typedefs = list(tdGrouping, "typedefs");
            typedefs.setEncodedKeys("ident");
            typedefs.addMeta(new Uses("def-header"));
            typedefs.addMeta(new Uses("type"));
        }
        {
            Grouping clluc = grouping(yang, "containers-lists-leafs-uses-choice");
            MetaList cllucList = list(clluc, "definitions");
            cllucList.setEncodedKeys("ident");

            Choice body = choice(cllucList, "body-stmt");

            ChoiceCase case1 = ccase(body, "container");
            Container container = contr(case1, "container");
            container.addMeta(new Uses("def-header"));
            container.addMeta(new Uses("groupings-typedefs"));
            container.addMeta(new Uses("containers-lists-leafs-uses-choice"));

            ChoiceCase case2 = ccase(body, "list");
            Container list = contr(case2, "list");
            list.addMeta(new Uses("def-header"));
            leafList(list, "keys", "string");
            list.addMeta(new Uses("groupings-typedefs"));
            list.addMeta(new Uses("containers-lists-leafs-uses-choice"));

            ChoiceCase case3 = ccase(body, "leaf");
            Container leaf = contr(case3, "leaf");
            leaf.addMeta(new Uses("def-header"));
            leaf(leaf, "config", "boolean");
            leaf(leaf, "mandatory", "boolean");
            leaf.addMeta(new Uses("type"));

            ChoiceCase case4 = ccase(body, "leaf-list");
            Container leafList = contr(case4, "leaf-list");
            leafList.addMeta(new Uses("def-header"));
            leaf(leafList, "config", "boolean");
            leaf(leafList, "mandatory", "boolean");
            leafList.addMeta(new Uses("type"));

            ChoiceCase case5 = ccase(body, "uses");
            Container uses = contr(case5, "uses");
            uses.addMeta(new Uses("def-header"));

            ChoiceCase case6 = ccase(body, "choice");
            Container choice = contr(case6, "choice");
            choice.addMeta(new Uses("def-header"));
            MetaList cases = list(choice, "cases");
            cases.setEncodedKeys("ident");
            leaf(cases, "ident", "string");
            cases.addMeta(new Uses("containers-lists-leafs-uses-choice"));
        }
        {
            Container module = contr(yang, "module");
            module.addMeta(new Uses("def-header"));
            leaf(module, "namespace", "string");
            leaf(module, "prefix", "string");
            Container rev = contr(module, "revision");
            leaf(rev, "rev-date", "string");
            leaf(rev, "description", "string");
            MetaList rpcs = list(module, "rpcs");
            rpcs.setEncodedKeys("ident");
            rpcs.addMeta(new Uses("def-header"));
            Container input = contr(rpcs, "input");
            input.addMeta(new Uses("groupings-typedefs"));
            input.addMeta(new Uses("containers-lists-leafs-uses-choice"));
            Container output = contr(rpcs, "output");
            output.addMeta(new Uses("groupings-typedefs"));
            output.addMeta(new Uses("containers-lists-leafs-uses-choice"));
            MetaList notifications = list(module, "notifications");
            notifications.setEncodedKeys("ident");
            notifications.addMeta(new Uses("def-header"));
            notifications.addMeta(new Uses("groupings-typedefs"));
            notifications.addMeta(new Uses("containers-lists-leafs-uses-choice"));
            module.addMeta(new Uses("groupings-typedefs"));
            module.addMeta(new Uses("containers-lists-leafs-uses-choice"));
        }
        YANG = yang;
    }

    static Grouping grouping(MetaCollection parent, String ident) {
        Grouping g = new Grouping(ident);
        parent.addMeta(g);
        return g;
    }

    static Choice choice(MetaCollection parent, String ident) {
        Choice c = new Choice(ident);
        parent.addMeta(c);
        return c;
    }

    static ChoiceCase ccase(Choice parent, String ident) {
        ChoiceCase c = new ChoiceCase(ident);
        parent.addMeta(c);
        return c;
    }

    static Container contr(MetaCollection parent, String ident) {
        Container child = new Container(ident);
        parent.addMeta(child);
        return child;
    }

    static MetaList list(MetaCollection parent, String ident) {
        MetaList child = new MetaList(ident);
        parent.addMeta(child);
        return child;
    }

    static Leaf leaf(MetaCollection parent, String ident, String type) {
        Leaf child = new Leaf(ident);
        child.setDataType(new DataType(type));
        parent.addMeta(child);
        return child;
    }

    static LeafList leafList(MetaCollection parent, String ident, String type) {
        LeafList child = new LeafList(ident);
        child.setDataType(new DataType(type));
        parent.addMeta(child);
        return child;
    }
}

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
            Grouping uses = grouping(yang, "uses");
            leaf(uses, "grouping", "string");
            MetaList refine = list(uses, "refine");
            leaf(refine, "default", "string");
        }
        {
            Grouping type = grouping(yang, "type");
            leaf(type, "grouping", "string");
            leaf(type, "range", "string");
            MetaList enumm = list(type, "enumeration");
            leaf(enumm, "range", "string");
        }
        {
            Grouping container = grouping(yang, "container");
            container.addMeta(new Uses("def-header"));
            container.addMeta(new Uses("uses"));
            container.addMeta(new Uses("definitions"));
        }
        {
            Grouping list = grouping(yang, "list");
            list.addMeta(new Uses("def-header"));
            list.addMeta(new Uses("uses"));
            list.addMeta(new Uses("definitions"));
        }
        {
            Grouping leaf = grouping(yang, "leaf");
            leaf.addMeta(new Uses("def-header"));
            leaf(leaf, "config", "string");
            leaf(leaf, "mandatory", "string");
            leaf.addMeta(new Uses("type"));
        }
        {
            Grouping leafList = grouping(yang, "leaf-list");
            leafList.addMeta(new Uses("def-header"));
            leaf(leafList, "config", "string");
            leaf(leafList, "mandatory", "string");
            leafList.addMeta(new Uses("type"));
        }
        {
            Grouping choice = grouping(yang, "choices");
            choice.addMeta(new Uses("def-header"));
            MetaList cases = list(choice, "cases");
            cases.addMeta(new Uses("definitions"));
        }
        {
            Grouping definitionsG = grouping(yang, "definitions");
            MetaList choices = list(definitionsG, "choices");
            choices.addMeta(new Uses("choices"));
            MetaList groupings = list(definitionsG, "groupings");
            groupings.addMeta(new Uses("groupings"));
            MetaList typedefs = list(definitionsG, "typedefs");
            typedefs.addMeta(new Uses("def-header"));
            typedefs.addMeta(new Uses("type"));
            MetaList definitions = list(definitionsG, "definitions");
            Choice body = choice(definitions, "body-stmt");
            definitions.addMeta(body);
            ChoiceCase case1 = ccase(body, "container");
            case1.addMeta(new Uses("container"));
            ChoiceCase case2 = ccase(body, "list");
            case2.addMeta(new Uses("list"));
            ChoiceCase case3 = ccase(body, "leaf");
            case3.addMeta(new Uses("leaf"));
            ChoiceCase case4 = ccase(body, "leaf-list");
            case4.addMeta(new Uses("leaf-list"));
            ChoiceCase case5 = ccase(body, "uses");
            case5.addMeta(new Uses("uses"));
        }
        {
            Grouping groupingsG = grouping(yang, "groupings");
            MetaList groupings = list(groupingsG, "groupings");
            groupings.addMeta(new Uses("def-header"));
            groupings.addMeta(new Uses("definitions"));
            groupings.addMeta(new Uses("uses"));
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
            rpcs.addMeta(new Uses("def-header"));
            Container input = contr(rpcs, "input");
            input.addMeta(new Uses("definitions"));
            Container output = contr(rpcs, "output");
            output.addMeta(new Uses("definitions"));
            MetaList notifications = list(module, "notifications");
            notifications.addMeta(new Uses("def-header"));
            notifications.addMeta(new Uses("definitions"));
            module.addMeta(new Uses("definitions"));
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
        child.setType(type);
        parent.addMeta(child);
        return child;
    }
}

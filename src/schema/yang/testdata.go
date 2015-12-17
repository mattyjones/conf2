package yang

var TestDataRomancingTheStone = `
module rtstone {
    namespace "rtstone";
    prefix "rtstone";
    revision 0000-00-00 {
        description "";
    }

    grouping position {
        leaf x {
            type int32;
        }
        leaf y {
            type int32;
        }
        leaf z {
            type int32;
        }
    }

    grouping team {
        leaf-list members {
            type string;
        }
        container spawn-point {
            uses position;
        }
        container base-position {
            uses position;
        }
    }

    container game {
        list teams {
            key "color";
            max-elements 4;
            leaf color {
                type string;
            }
            container team {
                uses team;
            }
        }
        leaf base-radius {
            type int32;
        }
        list leaderboard {
            container entry {
                leaf team {
                    type string;
                }
                leaf time {
                    type int32;
                }
            }
        }
        leaf time-limit {
            type int32;
        }
        list initial-inventory {
            container item {
                leaf item {
                    type int32;
                }
                leaf amount {
                    type int32;
                }
            }
        }
        list wardrobe {
            leaf item {
                type int32;
            }
        }
    }

	anyxml credits;

    rpc start-game {
        input {
            leaf seconds-from-now {
                type int32;
            }
        }
    }
    rpc end-game {
        input {
            leaf seconds-from-now {
                type int32;
            }
        }
    }
    rpc restart-game {
        input {
            leaf seconds-from-now {
                type int32;
            }
        }
    }
}`

var TestDataSimpleYang = `
module turing-machine {

  namespace "http://example.net/turing-machine";

  prefix "tm";

  description
    "Data model for the Turing Machine.";

  revision 2013-12-27 {
    description
      "Initial revision.";
  }

  /* Typedefs */

  typedef tape-symbol {
    description
      "Type of symbols appearing in tape cells.

       A blank is represented as an empty string where necessary.";
    type string {
      length "0..1";
    }
  }

  typedef cell-index {
    description
      "Type for indexing tape cells.";
    type int64;
  }

  typedef state-index {
    description
      "Type for indexing states of the control unit.";
    type uint16;
  }

  typedef head-dir {
    type enumeration {
      enum left;
      enum right;
    }
    default "right";
    description
      "Possible directions for moving the read/write head, one cell
       to the left or right (default).";
  }

  /* Groupings */

  grouping tape-cells {
    description
      "The tape of the Turing Machine is represented as a sparse
       array.";
    list cell {
      description
        "List of non-blank cells.";
      key "coord";
      leaf coord {
        type cell-index;
        description
          "Coordinate (index) of the tape cell.";
      }
      leaf symbol {
        type tape-symbol {
          length "1";
        }
        description
          "Symbol appearing in the tape cell.

           Blank (empty string) is not allowed here because the
           'cell' list only contains non-blank cells.";
      }
    }
  }

  /* State data and Configuration */

  container turing-machine {
    description
      "State data and configuration of a Turing Machine.";
    leaf state {
      config "false";
      mandatory "true";
      type state-index;
      description
        "Current state of the control unit.

         The initial state is 0.";
    }
    leaf head-position {
      config "false";
      mandatory "true";
      type cell-index;
      description
        "Position of tape read/write head.";
    }
    container tape {
      description
        "The contents of the tape.";
      config "false";
      uses tape-cells;
      action rewind {
        description "be kind";
        input {
          leaf position {
            type int32;
          }
        }
        output {
          leaf estimatedTime {
            type int32;
          }
        }
      }
    }
    container transition-function {
      description
        "The Turing Machine is configured by specifying the
         transition function.";
      list delta {
        description
          "The list of transition rules.";
        key "label";
        unique "input/state input/symbol";
        leaf label {
          type string;
          description
            "An arbitrary label of the transition rule.";
        }
        container input {
          description
            "Output values of the transition rule.";
          leaf state {
            type state-index;
            description
              "New state of the control unit. If this leaf is not
               present, the state doesn't change.";
          }
          leaf symbol {
            type tape-symbol;
            description
              "Symbol to be written to the tape cell. If this leaf is
               not present, the symbol doesn't change.";
          }
          leaf head-move {
            type head-dir;
            description
              "Move the head one cell to the left or right";
          }
        }
      }
    }
  }

  /* RPCs */

  rpc initialize {
    description
      "Initialize the Turing Machine as follows:

       1. Put the control unit into the initial state (0).

       2. Move the read/write head to the tape cell with coordinate
          zero.

       3. Write the string from the 'tape-content' input parameter to
          the tape, character by character, starting at cell 0. The
          tape is othewise empty.";
    input {
      leaf tape-content {
        type string;
        default "";
        description
          "The string with which the tape shall be initialized. The
           leftmost symbol will be at tape coordinate 0.";
      }
    }
  }

  rpc run {
    description
      "Start the Turing Machine operation.";
  }

  /* Notifications */

  notification halted {
    description
      "The Turing Machine has halted. This means that there is no
       transition rule for the current state and tape symbol.";
    leaf state {
      mandatory "true";
      type state-index;
    }
  }
}`

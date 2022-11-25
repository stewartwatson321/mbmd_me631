package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
        Register("IEM3000", NewIEM3000Producer)
}

type IEM3000Producer struct {
        Opcodes
}

func NewIEM3000Producer() Producer {
        /***
         * https://download.schneider-electric.com/files?p_enDocType=User+guide&p_File_Name=DOCA0005DE-12.pdf&p_Doc_Ref=DOCA0005DE#page49
         */
        ops := Opcodes{
                VoltageL1: 0x0863,
                VoltageL2: 0x0865,
                VoltageL3: 0x0867,
                Voltage:   0x0869,

                CurrentL1: 0x085B,
                CurrentL2: 0x085D,
                CurrentL3: 0x085F,
                Current:   0x0861,

                PowerL1: 0x086B,
                PowerL2: 0x086D,
                PowerL3: 0x086F,
                Power:   0x0871,

                ReactivePower: 0x0879,
                ApparentPower: 0x0881,

                // PowerFactor: 0x0C0B,
                Frequency: 0x07E6,

                Import:   0x0BC4,
                ImportL1: 0x0BB8,
                ImportL2: 0x0BBC,
                ImportL3: 0x0BC0,
                Export:   0x0BD4,

                ReactiveImport: 0x0BE4,
                ReactiveExport: 0x0BF4,
        }
        return &IEM3000Producer{Opcodes: ops}
}

// Description implements Producer interface
func (p *IEM3000Producer) Description() string {
        return "Schneider Electric iEM3000 series"
}

func (p *IEM3000Producer) snipFloat32(iec Measurement, scaler ...float64) Operation {
        snip := Operation{
                FuncCode:  ReadHoldingReg,
                OpCode:    p.Opcodes[iec],
                ReadLen:   2,
                IEC61850:  iec,
                Transform: RTUIeee754ToFloat64,
        }

        if len(scaler) > 0 {
                snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
        }

        return snip
}

func (p *IEM3000Producer) snipInt64(iec Measurement, scaler ...float64) Operation {
        snip := Operation{
                FuncCode:  ReadHoldingReg,
                OpCode:    p.Opcodes[iec],
                ReadLen:   4,
                IEC61850:  iec,
                Transform: RTUInt64ToFloat64,
        }

        if len(scaler) > 0 {
                snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
        }

        return snip
}

// Probe implements Producer interface
func (p *IEM3000Producer) Probe() Operation {
        return p.snipFloat32(VoltageL1)
}

// Produce implements Producer interface
func (p *IEM3000Producer) Produce() (res []Operation) {
        for op := range p.Opcodes {
                switch op {
                case PowerL1, PowerL2, PowerL3, Power, ReactivePower, ApparentPower:
                        res = append(res, p.snipFloat32(op, 0.001))
                case Import, ImportL1, ImportL2, ImportL3, Export, ReactiveImport, ReactiveExport:
                        res = append(res, p.snipInt64(op, 1000))
                default:
                        res = append(res, p.snipFloat32(op))
                }
        }

        return res
}

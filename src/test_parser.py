import json
import unittest

from google.protobuf.json_format import MessageToDict

from data import load
from parser import parse


class TestParser(unittest.TestCase):
    def test_jan_2020(self):
        self.maxDiff = None
        # Json com a saida esperada
        with open('src/output_test/expected/expected_01_2020.json', 'r') as fp:
            expected = json.load(fp)

        files = ['src/output_test/sheets/membros-ativos-contracheque-01-2020.csv',
                 'src/output_test/sheets/membros-ativos-verbas-indenizatorias-01-2020.csv']
                 
        dados = load(files, '2020', '01', 'src/output_test')
        result_data = parse(dados, 'mpro/01/2020', '01', '2020')
        # Converto o resultado do parser, em dict
        result_to_dict = MessageToDict(result_data)
        # with open('src/output_test/expected/expected_01_2020.json', 'w') as fp:
        #     json.dump(result_to_dict,fp,indent=4, ensure_ascii=False)
        self.assertEqual(expected, result_to_dict)


if __name__ == '__main__':
    unittest.main()